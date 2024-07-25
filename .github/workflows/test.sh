


version="v1.2.3"
IFS='.' read -ra version_arr <<< "${version}"

architectures=("amd64" "arm64" "arm" "386")
versions=("latest" "${version_arr[0]}.${version_arr[1]}" "${version}")

for arch in "${architectures[@]}"; do
for version in "${versions[@]}"; do
    IMAGE="${docker_registry}/tofuutils/tenv:${version}-${arch}"
    echo "Pushing ${IMAGE} ..."
    #docker push ${IMAGE}
    if [ ${?} -ne 0 ]; then
    echo "Failed to push ${IMAGE}"
    exit 1
    fi
done
done

echo "All images pushed successfully to ${docker_registry}!"
