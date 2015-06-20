#require './tippler'
require './tippler_web'
#require './lib/tippler/websocket'
#require 'faye/websocket'
#Faye::WebSocket.load_adapter('thin')
#Tippler::Websocket.init
#use Tippler::Websocket
run TipplerWeb
