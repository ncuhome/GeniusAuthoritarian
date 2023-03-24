import {MutableRefObject, useRef, useState} from "react";

interface RefState<T,> extends MutableRefObject<T> {
    state:T
}

export const useRefState= <V,>(initialValue:V):[RefState<V>,(v:V)=>void]=>{
    const [state,setState]=useState(initialValue)
    const ref=useRef(initialValue) as RefState<V>
    ref.state=state
    return [ref,(v:V)=>{
        ref.current=v
        setState(v)
    }]
}
