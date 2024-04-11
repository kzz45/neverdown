
function metadata_check(metadata: any) {
    let result = []
    if (!metadata.name) {
        result.push({
            key: 'metadata.name', tips: 'metadata.name is require!'
        })
    }
    return result
}
function spec_check(spec: any) {
    let result = []

    // if(spec.replicas !== 0 && !spec.replicas){
    //     result.push({
    //         key: 'spec.replicas', tips: 'spec.replicas is require!'
    //     })
    // }
    // if(spec.replicas < 0){
    //     result.push({
    //         key: 'spec.replicas', tips: 'spec.replicas need to be greater than or equal to 0'
    //     })
    // }
    return result
}

export default {
    metadata_check, spec_check
}
