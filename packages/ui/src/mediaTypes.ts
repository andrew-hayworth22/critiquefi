export type MediaTypeKey =
    "FILM" | "SHOW" | "SEASON" | "EPISODE" | "BOOK" | "GAME" | "ALBUM" | "SONG";

export interface MediaTypeInfo {
    key: MediaTypeKey;
    label: string;
    label_plural: string;
    count: string;
    isParent: boolean;
}

export const MEDIA_TYPES: MediaTypeInfo[] = [
    {
        key: "FILM",
        label: "Film",
        label_plural: "Films",
        count: "4.2M",
        isParent: true,
    },
    {
        key: "SHOW",
        label: "Show",
        label_plural: "Shows",
        count: "2.6M",
        isParent: true,
    },
    {
        key: "SEASON",
        label: "Season",
        label_plural: "Seasons",
        count: "",
        isParent: false,
    },
    {
        key: "EPISODE",
        label: "Episode",
        label_plural: "Episodes",
        count: "",
        isParent: false,
    },
    {
        key: "BOOK",
        label: "Book",
        label_plural: "Books",
        count: "3.1M",
        isParent: true,
    },
    {
        key: "GAME",
        label: "Game",
        label_plural: "Games",
        count: "1.8M",
        isParent: true,
    },
    {
        key: "ALBUM",
        label: "Album",
        label_plural: "Albums",
        count: "5.4M",
        isParent: true,
    },
    {
        key: "SONG",
        label: "Song",
        label_plural: "Songs",
        count: "",
        isParent: false,
    },
];

export const PARENT_MEDIA_TYPES: MediaTypeInfo[] = MEDIA_TYPES.filter(
    (m) => m.isParent,
);

export const MEDIA_TYPE_MAP: Record<MediaTypeKey, MediaTypeInfo> =
    Object.fromEntries(MEDIA_TYPES.map((type) => [type.key, type])) as Record<
        MediaTypeKey,
        MediaTypeInfo
    >;

export interface MediaTypeClasses {
    dot: string;
    card: string;
    label: string;
}

export const MEDIA_TYPE_CLASSES: Record<MediaTypeKey, MediaTypeClasses> = {
    FILM: {
        dot: 'bg-film',
        card: 'border-film/40 bg-linear-to-b from-film/16 to-film/3',
        label: 'text-film'
    },
    SHOW: {
        dot: 'bg-show',
        card: 'border-show/40 bg-linear-to-b from-show/16 to-show/3',
        label: 'text-show'
    },
    SEASON: {
        dot: 'bg-season',
        card: 'border-season/40 bg-linear-to-b from-season/16 to-season/3',
        label: 'text-season'
    },
    EPISODE: {
        dot: 'bg-episode',
        card: 'border-episode/40 bg-linear-to-b from-episode/16 to-episode/3',
        label: 'text-episode'
    },
    BOOK: {
        dot: 'bg-book',
        card: 'border-book/40 bg-linear-to-b from-book/16 to-book/3',
        label: 'text-book'
    },
    GAME: {
        dot: 'bg-game',
        card: 'border-game/40 bg-linear-to-b from-game/16 to-game/3',
        label: 'text-game'
    },
    ALBUM: {
        dot: 'bg-album',
        card: 'border-album/40 bg-linear-to-b from-album/16 to-album/3',
        label: 'text-album'
    },
    SONG: {
        dot: 'bg-song',
        card: 'border-song/40 bg-linear-to-b from-song/16 to-song/3',
        label: 'text-song'
    }
};
