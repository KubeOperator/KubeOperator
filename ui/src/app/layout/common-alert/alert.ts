export enum AlertLevels {
    SUCCESS, ERROR
}

export class Alert {
    msg: string;
    level: AlertLevels;

    constructor(msg: string, level: AlertLevels) {
        this.msg = msg;
        this.level = level;
    }
}
