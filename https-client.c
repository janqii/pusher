/*
  This is an example of how to hook up evhttp with bufferevent_ssl

  It just GETs an https URL given on the command-line and prints the response
  body to stdout.

  Actually, it also accepts plain http URLs to make it easy to compare http vs
  https code paths.

  Loosely based on le-proxy.c.
 */

// Get rid of OSX 10.7 and greater deprecation warnings.
#if defined(__APPLE__) && defined(__clang__)
#pragma clang diagnostic ignored "-Wdeprecated-declarations"
#endif

#include <stdio.h>
#include <assert.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>

#include <sys/socket.h>
#include <netinet/in.h>

#include <event2/bufferevent.h>
#include <event2/buffer.h>
#include <event2/listener.h>
#include <event2/util.h>
#include <event2/http.h>

static struct event_base *base;
static int ignore_cert = 0;

static void
http_request_done(struct evhttp_request *req, void *ctx)
{
}

static void
syntax(void)
{
	fputs("Syntax:\n", stderr);
	fputs("   https-client -url <https-url> [-data data-file.bin] [-ignore-cert] [-retries num]\n", stderr);
	fputs("Example:\n", stderr);
	fputs("   https-client -url https://ip.appspot.com/\n", stderr);

	exit(1);
}

static void
die(const char *msg)
{
	fputs(msg, stderr);
	exit(1);
}


int
main(int argc, char **argv)
{
	int r;

	struct evhttp_uri *http_uri;
	const char *url = NULL, *data_file = NULL;
	const char *scheme, *host, *path, *query;
	char uri[256];
	int port;
	int retries = 0;

	struct bufferevent *bev;
	struct evhttp_connection *evcon;
	struct evhttp_request *req;
	struct evkeyvalq *output_headers;
	struct evbuffer * output_buffer;

	int i;

	for (i = 1; i < argc; i++) {
		if (!strcmp("-url", argv[i])) {
			if (i < argc - 1) {
				url = argv[i + 1];
			} else {
				syntax();
			}
		} else if (!strcmp("-ignore-cert", argv[i])) {
			ignore_cert = 1;
		} else if (!strcmp("-data", argv[i])) {
			if (i < argc - 1) {
				data_file = argv[i + 1];
			} else {
				syntax();
			}
		} else if (!strcmp("-retries", argv[i])) {
			if (i < argc - 1) {
				retries = atoi(argv[i + 1]);
			} else {
				syntax();
			}
		} else if (!strcmp("-help", argv[i])) {
			syntax();
		}
	}

	if (!url) {
		syntax();
	}

	http_uri = evhttp_uri_parse(url);
	if (http_uri == NULL) {
		die("malformed url");
	}

	scheme = evhttp_uri_get_scheme(http_uri);
	if (scheme == NULL || (strcasecmp(scheme, "https") != 0 &&
	                       strcasecmp(scheme, "http") != 0)) {
		die("url must be http or https");
	}

	host = evhttp_uri_get_host(http_uri);
	if (host == NULL) {
		die("url must have a host");
	}

	port = evhttp_uri_get_port(http_uri);
	if (port == -1) {
		port = (strcasecmp(scheme, "http") == 0) ? 80 : 443;
	}

	path = evhttp_uri_get_path(http_uri);
	if (path == NULL) {
		path = "/";
	}

	query = evhttp_uri_get_query(http_uri);
	if (query == NULL) {
		snprintf(uri, sizeof(uri) - 1, "%s", path);
	} else {
		snprintf(uri, sizeof(uri) - 1, "%s?%s", path, query);
	}
	uri[sizeof(uri) - 1] = '\0';

	// Create event base
	base = event_base_new();
	if (!base) {
		perror("event_base_new()");
		return 1;
	}

    bev = bufferevent_socket_new(base, -1, BEV_OPT_CLOSE_ON_FREE);

    if (bev == NULL) {
        fprintf(stderr, "bufferevent_openssl_socket_new() failed\n");
        return 1;
    }


	// For simplicity, we let DNS resolution block. Everything else should be
	// asynchronous though.
	evcon = evhttp_connection_base_bufferevent_new(base, NULL, bev,
		host, port);
	if (evcon == NULL) {
		fprintf(stderr, "evhttp_connection_base_bufferevent_new() failed\n");
		return 1;
	}

	if (retries > 0) {
		evhttp_connection_set_retries(evcon, retries);
	}

	// Fire off the request
	req = evhttp_request_new(http_request_done, bev);
	if (req == NULL) {
		fprintf(stderr, "evhttp_request_new() failed\n");
		return 1;
	}

	output_headers = evhttp_request_get_output_headers(req);
	evhttp_add_header(output_headers, "Host", host);
	evhttp_add_header(output_headers, "Connection", "close");

	if (data_file) {
		/* NOTE: In production code, you'd probably want to use
		 * evbuffer_add_file() or evbuffer_add_file_segment(), to
		 * avoid needless copying. */
		FILE * f = fopen(data_file, "rb");
		char buf[1024];
		size_t s;
		size_t bytes = 0;

		if (!f) {
			syntax();
		}

		output_buffer = evhttp_request_get_output_buffer(req);
		while ((s = fread(buf, 1, sizeof(buf), f)) > 0) {
			evbuffer_add(output_buffer, buf, s);
			bytes += s;
		}
		evutil_snprintf(buf, sizeof(buf)-1, "%lu", (unsigned long)bytes);
		evhttp_add_header(output_headers, "Content-Length", buf);
		fclose(f);
	}

	r = evhttp_make_request(evcon, req, data_file ? EVHTTP_REQ_POST : EVHTTP_REQ_GET, uri);
	if (r != 0) {
		fprintf(stderr, "evhttp_make_request() failed\n");
		return 1;
	}

	event_base_dispatch(base);

	evhttp_connection_free(evcon);
	event_base_free(base);

	return 0;
}
