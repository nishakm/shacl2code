#! /usr/bin/env python3
#
# Copyright (c) 2024 Joshua Watt
#
# SPDX-License-Identifier: MIT

import argparse
import json
import sys
import urllib.request
from pathlib import Path

from . import Model
from .const import TEMPLATE_DIR
from .version import VERSION
from .lang import LANGUAGES


def get_langs():
    langs = []
    for child in TEMPLATE_DIR.iterdir():
        if child.suffixes and child.suffixes[-1] == ".j2":
            langs.append(child.stem)
    return langs


def main():
    def handle_generate(args):
        if "://" in args.input:
            with urllib.request.urlopen(args.input) as url:
                model_data = json.load(url)
        elif args.input == "-":
            model_data = json.load(sys.stdin)
        else:
            with Path(args.input).open("r") as f:
                model_data = json.load(f)

        m = Model(model_data)

        render = args.lang()
        render.output(args, m)
        return 0

    def handle_list(args):
        for k, v in LANGUAGES.items():
            if args.short:
                print(k)
            else:
                print(f"{k}: {v.HELP}")

    parser = argparse.ArgumentParser(
        description=f"Convert JSON-LD model to python. Version {VERSION}"
    )
    command_subparser = parser.add_subparsers(
        title="command",
        description="Command to execute",
        required=True,
    )
    generate_parser = command_subparser.add_parser(
        "generate",
        help="Generate language bindings",
    )
    generate_parser.add_argument(
        "--input",
        "-i",
        help="Input JSON-LD model (path, URL, or '-')",
        required=True,
    )
    generate_parser.set_defaults(func=handle_generate)

    lang_subparser = generate_parser.add_subparsers(
        title="language",
        description="Language to generate",
        required=True,
    )
    for k, v in LANGUAGES.items():
        p = lang_subparser.add_parser(k, help=v.HELP)
        v.get_arguments(p)
        p.set_defaults(lang=v)

    list_parser = command_subparser.add_parser("list", help="List languages")
    list_parser.add_argument(
        "--short",
        "-s",
        action="store_true",
        help="Only list languages without descriptions",
    )
    list_parser.set_defaults(func=handle_list)

    args = parser.parse_args()

    return args.func(args)