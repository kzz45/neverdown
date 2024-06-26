import { createI18n } from "vue-i18n"
export function loadLanguages() {
  const context = import.meta.globEager("./languages/*.ts");

  const languages: any = {};

  let langs = Object.keys(context);
  for (let key of langs) {
    if (key === "./index.ts") return;
    let lang = context[key].lang;
    let name = key.replace(/(\.\/languages\/|\.ts)/g, '');
    // try {
    //     if (name === "en") console.log('?????')
    //     // const elLang = await import(`element-plus/lib/locale/lang/en`) as  AnyObject;
    //     let elLang = await import(`element-plus/lib/locale/lang/${name}`) as AnyObject;
    //     lang = Object.assign(lang, {el: elLang.deafault.default.el})
    // } catch (error) {}
    languages[name] = lang
  }

  return languages
}

export const i18n = createI18n({
  legacy: false,
  locale: 'zh-cn',
  fallbackLocale: 'zh-cn',
  messages: loadLanguages(),
})

export function setLanguage(locale: string) {
  i18n.global.locale = locale
}

export const i18nt = i18n.global.t

export function getLanguage() {
  const language = (navigator.language || navigator.browserLanguage).toLowerCase()
  console.log('language', language)
  if (language === 'zh-cn') {
    return 'zh-cn'
  }
  return 'en'
}