#!/usr/bin/env python3

"""
    Performs initial configuration of Wikifeat services
    Note: Requires python3
"""

# from configparser import ConfigParser
from argparse import ArgumentParser
import util
from libs import configobj


user_config_template = '../users/config.ini.example'
user_config_file = '../users/config.ini'
user_node_id = 'us1'
user_port = 4100

wiki_config_template = '../wikis/config.ini.example'
wiki_config_file = '../wikis/config.ini'
wiki_node_id = 'ws1'
wiki_port = 4110

notifications_config_template = '../notifications/config.ini.example'
notifications_config_file = '../notifications/config.ini'
notifications_node_id = 'ns1'
notifications_port = 4120

frontend_config_template = '../frontend/config.ini.example'
frontend_config_file = '../frontend/config.ini'
frontend_node_id = 'fe1'
frontend_port = 8081
frontend_index_template = '../frontend/index.html.template'
frontend_index_file = '../frontend/index.html'
frontend_plugin_template = '../frontend/plugins/plugins.ini.example'
frontend_plugin_file = '../frontend/plugins/plugins.ini'


def config_user_service(common_params, db_params):
    print("Configuring users service...")
    try:
        config = configobj.ConfigObj(user_config_template, file_error=True)
    except IOError:
        return False
    # Configure the Service section
    config['Service']['domainName'] = common_params['domainName']
    config['Service']['nodeId'] = user_node_id
    config['Service']['port'] = str(user_port)
    # Now the database section
    config_database(config, db_params)
    # Now, write the config file
    with open(user_config_file, 'w') as out_file:
        config.write(out_file)
    return True


def config_wiki_service(common_params, db_params):
    print("Configuring wikis service...")
    try:
        config = configobj.ConfigObj(wiki_config_template, file_error=True)
    except IOError:
        return False
    config_service(config, common_params['domainName'],
                   wiki_node_id, wiki_port)
    config_database(config, db_params)
    with open(wiki_config_file, 'w') as out_file:
        config.write(out_file)
    return True


def config_notifications_service(common_params, db_params):
    print("Configuring notifications service...")
    try:
        config = configobj.ConfigObj(
            notifications_config_template, file_error=True)
    except IOError:
        return False
    config_service(config, common_params['domainName'],
                   notifications_node_id, notifications_port)
    config_database(config, db_params)
    main_site_url = "http://%s" % common_params['domainName']
    if frontend_port != 80 and frontend_port != 443:
        main_site_url += ":%s" % str(frontend_port)
    config['Notifications']['mainSiteUrl'] = main_site_url
    with open(notifications_config_file, 'w') as out_file:
        config.write(out_file)
    return True


def config_frontend_service(common_params, db_params):
    import shutil
    print("Configuring frontend service...")
    try:
        config = configobj.ConfigObj(
            frontend_config_template, file_error=True
        )
    except IOError:
        return False
    config_service(config, common_params['domainName'],
                   frontend_node_id, frontend_port)
    config_database(config, db_params)
    with open(frontend_config_file, 'w') as out_file:
        config.write(out_file)
    shutil.copyfile(frontend_index_template, frontend_index_file)
    shutil.copyfile(frontend_plugin_template, frontend_plugin_file)
    return True


def config_service(config, domain_name, node_id, port):
    config['Service']['domainName'] = domain_name
    config['Service']['nodeId'] = node_id
    config['Service']['port'] = str(port)


def config_database(config, db_params):
    config['Database']['dbAddr'] = db_params.host
    config['Database']['dbPort'] = db_params.port
    config['Database']['dbAdminUser'] = db_params.adminuser
    config['Database']['dbAdminPassword'] = db_params.adminpass


def config_all(common_params, db_params):
    print("Configuring services...")
    config_user_service(common_params, db_params)
    config_wiki_service(common_params, db_params)
    config_frontend_service(common_params, db_params)
    config_notifications_service(common_params, db_params)
    print("Configuration complete")


def main(domain_name, db_params):
    common_params = dict()
    common_params['domainName'] = domain_name
    config_all(common_params, db_params)


if __name__ == "__main__":
    parser = ArgumentParser()
    util.add_couch_params(parser)
    parser.add_argument('--domain_name', dest='domain_name',
                        help='host domain name')
    parser.set_defaults(domain_name='localhost')
    args = parser.parse_args()
    if args.adminuser is None:
        args.adminuser = input("Enter CouchDB admin username: ")
    if args.adminpass is None:
        args.adminpass = input("Enter CouchDB admin password: ")
    couch_params = util.CouchParameters(args)
    main(args.domain_name, couch_params)



