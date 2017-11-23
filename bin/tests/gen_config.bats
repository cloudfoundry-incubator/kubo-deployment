#!/usr/bin/env bats

setup(){
    tmp_dir="$(mktemp -d)"
    echo $tmp_dir
    #cp -r ../src/kubo-deployment-tests/resources/environments/test_gcp_with_creds "${tmp_dir}"
    #sed -n 's/iaas: gcp/iaas: foobar/g' ${tmp_dir}/director.yml > ${tmp_dir}/director.yml
}

teardown(){
    echo "Teardown"
    #rm -rf "${tmp_dir}"
}

@test "print usage for generate_env_config" {
    run bash generate_env_config
    [ "$status" -eq 1 ]
    [ "${lines[3]}" = "Usage: generate_env_config PATH NAME IAAS" ]
}

@test "incorrect number of parameters passed to generate_env_config" {
    run bash generate_env_config /tmp/abcd foo
    [ "$status" -eq 1 ]
    [ "${lines[3]}" = "Usage: generate_env_config PATH NAME IAAS" ]
}

@test "non-existent directory passed to generate_env_config" {
    run bash generate_env_config /tmp/abcd foo
    [ "$status" -eq 1 ]
    [ "${lines[3]}" = "Usage: generate_env_config PATH NAME IAAS" ]
}

@test "incorrect iaas passed to generate_env_config" {
    run bash generate_env_config $tmp_dir foo xyz
    [ "$status" -eq 1 ]
    [ "${lines[3]}" = "Usage: generate_env_config PATH NAME IAAS" ]
}

@test "success run of generate_env_config" {
    run bash generate_env_config $tmp_dir ci-service gcp
    [ -d "$tmp_dir/ci-service" ]
    [ "$status" -eq 0 ]

    run bash -c "ls -1 $tmp_dir/ci-service | wc -l"
    [ "$output" -eq 2 ]

    run grep "iaas: gcp" $tmp_dir/ci-service/director.yml
    [ "$output" = "iaas: gcp" ]
    [ "$status" -eq 0 ]
}