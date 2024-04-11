// import Base from "@/lib/ts/Base";
import { InjectionKey } from "vue"
import { loadModules, context, modules } from "./modules"
import { createStore, useStore as baseUseStore, Store, createLogger } from "vuex"
const IS_DEV = process.env.NODE_ENV == 'development'
export interface State {
    [key: string]: any
}

export const key: InjectionKey<Store<State>> = Symbol();

const store = createStore({
    modules,
    strict: IS_DEV,
    plugins: []
    // plugins: Base.IS_DEV ? [createLogger()] : [] //store log
});

export function useStore() {
    // return baseUseStore(key);
    return baseUseStore();
}

// 热重载
if (import.meta.hot) {
    import.meta.hot?.accept(context.id, () => {
        const { modules } = loadModules()
        store.hotUpdate({
            modules
        })
    })
}

export default store;