import {apiV1} from "@api/base";

export async function GetFeishuLoginUrl(target:string):Promise<string>{
    const {data:{data:{url}}} = await apiV1.get('public/login/feishu/link', {
        params: {
            callback: target
        }
    })
    return url
}
