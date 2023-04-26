import {FC} from "react";
import {useMount, useQuery} from "@hooks";
import {useNavigate} from "react-router-dom";
import {ThrowError} from "@util/nav";

import {OnLogin} from "@components";

import {DingTalkLogin} from "@api/v1/login";

export const DingTalk: FC = () => {
    const nav = useNavigate();
    const [code] = useQuery("authCode", "");
    const [callback] = useQuery("state", "");

    async function login() {
        try {
            const callbackUrl = await DingTalkLogin(code, callback);
            window.open(callbackUrl, "_self")
        } catch ({msg}) {
            ThrowError(nav, "登录失败", msg as string)
        }
    }

    useMount(() => {
        if (!code || !callback) {
            ThrowError(nav, "登录失败", "参数缺失");
            return;
        }
        login()
    });

    return <OnLogin/>;
}
