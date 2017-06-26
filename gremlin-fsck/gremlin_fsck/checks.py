from gremlin_python.process.graph_traversal import __
from gremlin_python.process.traversal import within


def check_vn_with_iip_without_vmi(g):
    """find instance-ip without any virtual-machine-interface
    """
    return g.V().hasLabel("virtual_network").not_(
        __.in_().hasLabel('virtual_machine_interface')
    ).in_().hasLabel("instance_ip")


def clean_vn_with_iip_without_vmi(iip):
    iip.delete()


def check_unused_rt(g):
    """find unused route-target
    """
    return g.V().hasLabel("route_target").not_(
        __.in_().hasLabel(within("routing_instance", "logical_router"))
    )


def clean_unused_rt(rt):
    rt.delete()


def check_iip_without_instance_ip_address(g):
    """find iip without any instance_ip_address property
    """
    return g.V().hasLabel("instance_ip").not_(
        __.has("instance_ip_address")
    )


def check_snat_without_lr(g):
    """find Snat SI without any logical-router
    """
    return g.V().hasLabel("service_template").has("name", "netns-snat-template") \
        .in_().hasLabel("service_instance").not_(__.in_().hasLabel("logical_router"))


def check_lbaas_without_lbpool(g):
    """find LBaaS SI without any loadbalancer-pool
    """
    return g.V().hasLabel("service_template") \
        .has("name", "haproxy-loadbalancer-template") \
        .in_().hasLabel("service_instance") \
        .not_(__.in_().hasLabel("loadbalancer_pool"))


def check_lbaas_without_vip(g):
    """find LBaaS SI without any virtual-ip
    """
    return g.V().hasLabel("service_instance") \
        .where(__.in_().hasLabel("loadbalancer_pool").not_(__.in_().hasLabel("virtual_ip")))


def check_ri_without_rt(g):
    """find routing-instance that doesn't have any route-target (that crashes schema)
    """
    return g.V().hasLabel("routing_instance") \
        .not_(__.has('fq_name', within("__default__", "__link_local__"))) \
        .not_(__.out().hasLabel("route_target"))