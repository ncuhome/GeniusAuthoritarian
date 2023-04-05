import axios from "axios";

export const BaseURL = `${location.protocol}://${location.host}/api/`
export const BaseUrlV1 = `${BaseURL}v1/`


const apiV1 = axios.create({
    baseURL: BaseUrlV1
})
apiV1.interceptors.response.use(undefined, (err:any):any=>{
    if(err.name==="CanceledError") {
        return new Promise(() => {})
    }
    if(!err)err={}
    if(!err.response)err.response={}
    if(!err.response.data)err.response.data={code:-1}
    return err
})


export {apiV1}
