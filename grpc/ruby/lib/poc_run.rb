require_relative 'services_pb'
require_relative 'services_services_pb'

stub = Jjogaegi::RunService::Stub.new('localhost:5000', :this_channel_is_insecure)

resp = stub.run(Jjogaegi::RunRequest.new(
    config: Jjogaegi::RunConfig.new(
      parser:    "list",
      formatter: "csv",
      options: {
         'debug' => 'true'
      }
    ),
    input: "한국 yyy".b,
))

puts resp.output