#!/usr/bin/python


import os

from pyln.testing.fixtures import *  # noqa: F403

test_binary = os.path.join(os.path.dirname(__file__), "go-lnmetrics")


def test_basic(node_factory):
    node = node_factory.get_node(
        options={
            "plugin": test_binary,
            "lnmetrics-urls": "https://api.lnmetrics.info/query",
            "lnmetrics-noproxy": True,
        }
    )
    # node.rpc.call("metric_one", {"start": "now"})
    # node.rpc.call("raw-local-score")
    # node.rpc.call("lnmetrics-force-update")
    # node.rpc.call("lnmetrics-info")
    # node.rpc.call("lnmetrics-clean")
    assert False
