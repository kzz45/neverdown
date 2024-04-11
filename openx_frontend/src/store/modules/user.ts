import { Module } from "vuex";
import { getLanguage } from '@/i18n'

interface StoreUser {
    //lang: string,
    namespec: any[],
    compareItems: any[]
}

const store: Module<StoreUser, unknown> = {
    namespaced: true,
    state() {
        return {
            // lang: getLanguage(),
            namespec: [],
            compareItems: []
        }
    },
    mutations: {
        setCompareItems(state: StoreUser, payload: any) {
            state.compareItems = payload;
        },
        setSpec(state: StoreUser, payload: any) {
            state.namespec = payload.namespec;
        },
    },
    actions: {
        setCompareItems(context, payload: any) {
            context.commit("setCompareItems", payload);
        },
        setSpec(context, payload: any) {
            context.commit("setSpec", payload);
        }
    },
    getters: {
        getCompareItems(state: StoreUser) {
            return state.compareItems
        },
        getSpec(state: StoreUser) {
            return state.namespec
        }
    }
}

export default store