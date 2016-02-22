#!/usr/bin/env python3

"""
    Performs initial configuration of Wikifeat services
    Note: Requires python3
"""

# from configparser import ConfigParser
from argparse import ArgumentParser
import util
import os
from libs import configobj

wikifeat_path = os.path.join(os.getcwd(), os.pardir)


def config_template():
    return os.path.realpath(wikifeat_path + '/config/config.ini.example')


def config_file():
    return os.path.realpath(wikifeat_path + '/config/config.ini')


def frontend_index_template():
    return os.path.realpath(wikifeat_path + '/frontend/index.html.template')


def frontend_index_file():
    return os.path.realpath(wikifeat_path + '/frontend/index.html')


def frontend_plugin_template():
    return os.path.realpath(wikifeat_path + '/frontend/plugins/plugins.ini.example')


def frontend_plugin_file():
    return os.path.realpath(wikifeat_path + '/frontend/plugins/plugins.ini')


def config_frontend_service():
    import shutil
    print("Configuring frontend service...")
    shutil.copyfile(frontend_index_template(), frontend_index_file())
    shutil.copyfile(frontend_plugin_template(), frontend_plugin_file())
    return True


def config_database(config, db_params):
    config['Database']['dbAddr'] = db_params.host
    config['Database']['dbPort'] = db_params.port
    config['Database']['dbAdminUser'] = db_params.adminuser
    config['Database']['dbAdminPassword'] = db_params.adminpass


def config_webapp(config):
    print("Configuring webapp...")
    config['Frontend']['webAppDir'] = wikifeat_path + "/frontend/web_app/app"
    config['Frontend']['pluginDir'] = wikifeat_path + "/frontend/plugins"

def config_all(common_params, db_params):
    print("Configuring wikifeat...")
    try:
        config = configobj.ConfigObj(
            config_template(), file_error=True
        )
    except IOError:
        print("Error reading file " + config_template())
        return False
    config_database(config, db_params)
    config_webapp(config)
    with open(config_file(), 'w') as out_file:
        config.write(out_file)
    config_frontend_service()
    print("Configuration complete")


def main(domain_name, db_params, install_dir):
    global wikifeat_path
    if install_dir is not None:
        wikifeat_path = install_dir
    common_params = dict()
    common_params['domainName'] = domain_name
    config_all(common_params, db_params)


if __name__ == "__main__":
    parser = ArgumentParser()
    util.add_couch_params(parser)
    parser.add_argument('--domain_name', dest='domain_name',
                        help='host domain name')
    parser.add_argument('--wikifeat_home', dest='wikifeat_home',
                        help='Wikifeat install directory')
    parser.set_defaults(domain_name='localhost')
    parser.set_defaults(wikifeat_home=os.curdir)
    args = parser.parse_args()
    if args.adminuser is None:
        args.adminuser = input("Enter CouchDB admin username: ")
    if args.adminpass is None:
        args.adminpass = input("Enter CouchDB admin password: ")
    couch_params = util.CouchParameters(args)
    main(args.domain_name, couch_params, args.wikifeat_home)



