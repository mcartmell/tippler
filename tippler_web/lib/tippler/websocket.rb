require 'byebug'
class Tippler
  class Websocket
    class << self
      attr_reader :tippler

      def init
        @tippler = Tippler.new
        EM.next_tick { @tippler.start_twitter_stream }
      end
    end

    def initialize(app)
      @app = app
    end

    def call(env)
      if Faye::WebSocket.websocket?(env)
        ws = Faye::WebSocket.new(env, ping: 2)

        ws.on :open do |event|
          p [:open, ws.object_id]
          self.class.tippler.client_connected(ws)
        end

        ws.on :close do |event|
          self.class.tippler.client_disconnected(ws)
          p [:close, event.code, event.reason]
          ws = nil
        end

        # Return async Rack response
        ws.rack_response

      else
        @app.call(env)
      end
    end
  end
end
