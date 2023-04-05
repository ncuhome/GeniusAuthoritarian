import {apiV1} from "@api/base";

export async function GetFeishuLoginUrl():Promise<string>{
    const {data:{data:{url}}} = await apiV1.get('public/login/feishu/link', {
        params: {
            callback: `${location.protocol}://${location.host}/feishu`
        }
    })
    return url
}
