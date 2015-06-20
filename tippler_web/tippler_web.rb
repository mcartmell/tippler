require 'sinatra/base'
require 'sinatra/asset_pipeline'
class TipplerWeb < Sinatra::Base
  set :assets_prefix, %w{assets vendor/assets}
  register Sinatra::AssetPipeline

  Areas = ["Ang Mo Kio", "Bedok", "Bishan", "Boon Lay", "Bukit Batok", "Bukit Merah", "Bukit Panjang", "Bukit Timah",
            "Central Business District", "Central Water Catchment", "Changi", "Changi Bay", "Choa Chu Kang", "Clementi",
            "Geylang", "Hougang", "Jurong East", "Jurong West", "Kallang", "Lim Chu Kang", "Mandai", "Marina East",
            "Marina South", "Marine Parade", "Museum", "Newton", "North-eastern Islands", "Novena", "Orchard", "Outram",
            "Pasir Ris", "Paya Lebar", "Pioneer", "Punggol", "Queenstown", "River Valley", "Rochor", "Seletar",
            "Sembawang", "Sengkang", "Serangoon", "Simpang", "Singapore River", "Southern Islands", "Straits View",
            "Sungei Kadut", "Tampines", "Tanglin", "Tengah", "Toa Payoh", "Tuas", "Western Islands", "Western Water
            Catchment", "Woodlands", "Yishun"]
  get '/' do
    @current_channel = 'All'
    @all_areas = Areas
    erb :index
  end
end

