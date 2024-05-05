const SwaggerParser = require("@apidevtools/swagger-parser");

async function main() {
  const startTime = new Date().getTime();
  let endTime;
  const data = await SwaggerParser.parse("data/test_openapi_spec.yaml");
  endTime = new Date().getTime();
  console.log("Time taken: " + (endTime - startTime) + "ms");
  console.log(data.openapi)
}

main()