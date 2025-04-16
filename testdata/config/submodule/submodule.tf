resource "random_pet" "server" {
  count = 3
  keepers = {
    foo = "null_data_source.count.${count.index}.id"
  }
}

data "null_data_source" "values" {
  inputs = {
    all_server_ids = join("", random_pet.server.*.id)
  }
}

data "null_data_source" "for-each" {
  for_each = toset( ["Todd", "James", "Alice", "Dottie"] )
  inputs = {
    foo = each.key
  }
}

resource "null_resource" "cluster" {
  count = 3
}

data "null_data_source" "count" {
  count = 3
  inputs = {
    bar = "null_resource.cluster.${count.index}.id"
  }
}

output "all_server_ids" {
  value = data.null_data_source.values.outputs["all_server_ids"]
}
