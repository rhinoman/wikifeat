"""
    Common utility functions for setup scripts
"""

import json, os
from argparse import ArgumentParser
import http.client


def parse_args():
    parser = ArgumentParser()
    add_couch_params(parser)
    parser.add_argument('-db', '--main_db',
                        dest='main_db',
                        help='Main Wikifeat database')
    parser.add_argument('-adb', '--avatar_db',
                        dest='avatar_db',
                        help='User Avatar database')
    parser.add_argument('--master_username', dest='m_uname',
                        help="Master User name")
    parser.add_argument('--master_password', dest='m_password',
                        help="Master User password")
    parser.add_argument('--master_firstname', dest='m_fn',
                        help="Master User first name")
    parser.add_argument('--master_lastname', dest='m_ln',
                        help="Master User last name")
    parser.add_argument('--skip_master', dest='skip_master', action='store_true',
                        help="Skip master user setup")
    parser.add_argument('--wikifeat_home', dest='wikifeat_home',
                        help='Installation directory for Wikifeat')

    parser.set_defaults(main_db='wikifeat_main_db')
    parser.set_defaults(avatar_db='user_avatars')
    parser.set_defaults(skip_master=False)
    parser.set_defaults(wikifeat_home=os.path.realpath(os.path.join(os.curdir, os.pardir)))

    args = parser.parse_args()

    if args.adminuser is None:
        args.adminuser = input("Enter CouchDB admin username: ")

    if args.adminpass is None:
        args.adminpass = input("Enter CouchDB admin password: ")

    return args


def add_couch_params(parser):
    parser.add_argument('couch_server', type=str,
                        help='CouchDB host')
    parser.add_argument('couch_port', type=int,
                        help='CouchDB port')
    parser.add_argument('-u', '--user', dest='adminuser',
                        help='CouchDB admin user')
    parser.add_argument('-p', '--password', dest='adminpass',
                        help='CouchDB admin password')
    # Note: your python must be compiled with SSL support
    parser.add_argument('--use_ssl', dest='use_ssl', action='store_true',
                        help="Use SSL to connect to CouchDB.  Your python must "
                             "have been compiled with SSL support!")

    parser.set_defaults(use_ssl=False)


class MasterUserParameters(object):
    def __init__(self, args=None):
        if args is not None:
            self.user = args.m_uname
            self.password = args.m_password
            self.firstname = args.m_fn
            self.lastname = args.m_ln
            self.skip_master = args.skip_master
        else:
            self.user = None
            self.password = None
            self.firstname = None
            self.lastname = None
            self.skip_master = False


class CouchParameters(object):
    def __init__(self, args=None):
        if args is not None:
            self.host = args.couch_server
            self.port = args.couch_port
            self.adminuser = args.adminuser
            self.adminpass = args.adminpass
            self.use_ssl = args.use_ssl
        else:
            self.host = "localhost"
            self.port = 5984
            self.use_ssl = False


def get_connection(use_ssl, couch_server, couch_port):
    if use_ssl:
        return http.client.HTTPSConnection(couch_server, couch_port)
    else:
        return http.client.HTTPConnection(couch_server, couch_port)


def load_json_file(filename):
    with open(filename) as json_file:
        data = json.load(json_file)
        return data


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


def decode_response(resp):
    return json.loads(resp.read().decode('utf-8'))
