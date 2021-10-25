#!/usr/bin/env python3

# LSST Data Management System
# Copyright 2014 LSST Corporation.
# 
# This product includes software developed by the
# LSST Project (http://www.lsst.org/).
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
# 
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
# 
# You should have received a copy of the LSST License Statement and 
# the GNU General Public License along with this program.  If not, 
# see <http://www.lsstcorp.org/LegalNotices/>.

"""
Create k8s Persistent Volumes and Persistent Volume Claims

@author Benjamin Roziere, IN2P3
@author Fabrice Jammes, IN2P3
"""

# -------------------------------
#  Imports of standard modules --
# -------------------------------
import argparse
import os.path
import sys
import yaml

def _build_yaml(data_path, pvc_name, hostname, instance, output_dir, template_dir):

    # yaml for persistent volume
    #
    tpl_fname = 'pv-qserv.tpl.yaml'

    yaml_tpl = os.path.join(template_dir, tpl_fname)
    with open(yaml_tpl, 'r') as f:
        yaml_data = yaml.load(f, Loader=yaml.SafeLoader)

    yaml_data['metadata']['name'] = "pv-{}".format(pvc_name)
    yaml_data['metadata']['labels']['pvc_name'] = pvc_name
    yaml_data['metadata']['labels']['instance'] = instance

    node_name = yaml_data['spec']['nodeAffinity']['required']['nodeSelectorTerms'][0]['matchExpressions'][0]['values']
    node_name[0] = hostname
    yaml_data['spec']['local']['path'] = data_path

    yaml_fname = "pv-{}.yaml".format(pvc_name)
    yaml_fname = os.path.join(output_dir, yaml_fname)
    with open( yaml_fname, "w") as f:
        f.write(yaml.dump(yaml_data, default_flow_style=False))

    # yaml for persistent volume claim
    #
    yaml_tpl = os.path.join(template_dir, 'pvc-qserv.tpl.yaml')
    with open(yaml_tpl, 'r') as f:
        yaml_data = yaml.load(f, yaml.SafeLoader)

    yaml_data['metadata']['name'] = "{}".format(pvc_name)
    yaml_data['metadata']['labels']['instance'] = instance
    yaml_data['spec']['selector']['matchLabels']['pvc_name'] = pvc_name

    yaml_fname = "pvc-{}.yaml".format(pvc_name)
    yaml_fname = os.path.join(output_dir, yaml_fname)
    with open( yaml_fname, "w") as f:
        f.write(yaml.dump(yaml_data, default_flow_style=False))

if __name__ == "__main__":
    try:

        cdir = os.path.dirname(os.path.realpath(__file__))
        manifest_dir = os.path.join(cdir, "manifests")

        parser = argparse.ArgumentParser(description="Create k8s Persistent Volumes and Claims")

        parser.add_argument('-p', '--path', dest='data_path',
                            required=True, metavar='<hostPath>',
                            help='Path on the host')
        parser.add_argument('-n', '--pvcname', dest='pvc_name',
                            required=True, metavar='<persistentVolumeClaimName>',
                            help='Name of the PersistentVolumeClaim')
        parser.add_argument('-H', '--hostname', dest='hostname',
                            required=False, metavar='<hostname>',
                            help='Hostname of the node')
        parser.add_argument('-t', '--templateDir', dest='template_dir',
                            default=manifest_dir,
                            required=False, metavar='<templateDir>',
                            help='yaml template directory')
        parser.add_argument('-o', '--outputDir', dest='output_dir',
                            required=True, metavar='<outputDir>',
                            help='Output directory for generated yaml files')
        parser.add_argument('-i', '--instance', dest='instance',
                            required=True, metavar='<instance>',
                            help='Name of qserv instance')

        args = parser.parse_args()

        _build_yaml(args.data_path, args.pvc_name, args.hostname, args.instance, args.output_dir, args.template_dir)

    except Exception as e:
        print(e)
        sys.exit(1)
