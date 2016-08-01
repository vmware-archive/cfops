require 'sinatra/base'

class TestApp < Sinatra::Base
  get "/" do
    "hello world"
  end
end
