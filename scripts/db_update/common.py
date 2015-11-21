__author__ = 'jcadam'

import json
import http.client
from argparse import ArgumentParser


def decode_response(resp):
        return json.loads(resp.read().decode('utf-8'))


def parse_args():
    # We need the CouchDB admin credentials
    # These can be provided as command line arguments or via prompt

    parser = ArgumentParser()
    # CouchDB URL is required
    parser.add_argument('couch_server', type=str)
    parser.add_argument('couch_port', type=int)
    parser.add_argument('-u', '--user', dest='adminuser',
                        help='CouchDB admin user')
    parser.add_argument('-p', '--password', dest='adminpass',
                        help='CouchDB admin password')
    parser.add_argument('-db', '--maindb', dest='main_db',
                        help='Main Wikifeat database')
    # Note: your python must be compiled with SSL support to use HTTPS
    parser.add_argument('--use_ssl', dest='use_ssl', action='store_true')
    parser.set_defaults(use_ssl=False)
    parser.set_defaults(main_db="wikifeat_main_db")

    args = parser.parse_args()

    if args.adminuser is None:
        args.adminuser = input("Enter CouchDB admin username: ")

    if args.adminpass is None:
        args.adminpass = input("Enter CouchDB admin password: ")

    return args


def get_connection(use_ssl, couch_server, couch_port):
    if use_ssl:
        conn = http.client.HTTPSConnection(couch_server, couch_port)
    else:
        conn = http.client.HTTPConnection(couch_server, couch_port)
    return conn


def get_credentials(user, password):
    from base64 import b64encode
    credentials = b64encode(bytes(user + ':' + password, 'utf-8')).decode('utf-8')
    return credentials


def get_headers(credentials):
    return {
        'Accept': 'application/json',
        'Authorization': 'Basic %s' % credentials
    }


def put_headers(credentials):
    p_headers = get_headers(credentials).copy()
    p_headers.update({'Content-Type': 'application/json'})
    return p_headers

