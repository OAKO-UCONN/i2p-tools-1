__doc__ = """
application layer filtering i2p client/server tunnels
"""


import logging

from pyi2ptunnel import tunnels

from twisted.internet import reactor
from twisted.internet.endpoints import clientFromString, quoteStringArgument

log = logging.getLogger("pyi2ptunnel.tunnel")


class TunnelFactory(object):
    """
    creator of app tunnels
    """

    def __init__(self, api="SAM", apiEndpoint="tcp:127.0.0.1:7656", socksPort=9050):
        self.api = api
        self.apiEndpoint = quoteStringArgument(apiEndpoint)
        self.socksPort = socksPort

    def endpoint(self, addr, port):
        """
        return the endpoint
        """
        if addr.endswith('.i2p'):
            ep = 'i2p:{}:api={}:apiEndpoint={}'.format(addr, self.api, self.apiEndpoint)
        else:
            ep = 'tor:host={}:port={}:socksPort={}'.format(addr, port, self.socksPort)
        return clientFromString(reactor, ep)
                                
                                    
        
    def createClient(self, type, **param):
        if type in tunnels.clients:
            return tunnels.clients[type](self.endpoint, **param)
        
    def createServer(self, type, **param):
        if type in tunnels.servers:
            return tunnels.servers[type](self.endpoint, **param)
