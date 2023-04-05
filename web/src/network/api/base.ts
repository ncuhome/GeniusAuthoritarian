import axios from "axios";

export const BaseURL = `/api/`
export const BaseUrlV1 = `${BaseURL}v1/`

const apiV1 = axios.create({
    baseURL: BaseUrlV1
})
apiV1.interceptors.response.use(undefined, (err:any)=>{
    switch (true) {
        case err.name==="CanceledError":
            break
        case !err||!err.response||!err.response.data:
            err.msg="网络异常，请检查网络设置"
            break
        default:
            err.msg=err.response.data.msg
    }
    return Promise.reject(err)
})


export {apiV1}
