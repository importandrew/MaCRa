import type en from './locales/en';

// PluralForm is a CLDR plural category i18next appends to a key ("games_few").
// English uses one and other; other languages add more (Russian: few, many).
type PluralForm = 'zero' | 'one' | 'two' | 'few' | 'many' | 'other';

// Translation is the shape of T with every text replaced by string: the same
// keys and nesting, any wording. A group may also carry plural forms the
// reference language doesn't need.
type Translation<T> = {
  [K in keyof T]: T[K] extends string ? string : Translation<T[K]>;
} & { [plural: `${string}_${PluralForm}`]: string };

// Messages is what every language file must provide: the shape of the
// English texts.
export type Messages = Translation<typeof en>;
