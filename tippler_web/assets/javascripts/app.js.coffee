#= require vendor
#= require_self
class Tippler
  init: ->
    @connect('all')
    @bindings()

  connect: (channel) ->
    if @es
      @es.close()
    @es = new WebSocket("ws://localhost:9292/" + channel)
    @es.onmessage = (event) ->
      d = JSON.parse(event.data)
      $('#tweets').append('<li><span class="time">[' + d.time + ']</span>&nbsp;<span class="author">&lt;' + d.user + '&gt;</span>&nbsp;' + d.msg + '</li>')
      el = $('.main-section')[0]
      el.scrollTop = el.scrollHeight
    $('#tweets').empty()

  bindings: ->
    tip = this
    $('[data-channel]').on 'click', ->
      $('.off-canvas-wrap').foundation('offcanvas', 'hide', 'move-right');
      chan = $(this).data 'channel'
      label = $(this).data 'channel-label'
      $('#current_area').html(label)
      tip.connect(chan)
$ ->
  t = new Tippler
  t.init()
