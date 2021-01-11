import {Component, OnInit} from '@angular/core';
import {ClusterLog} from "../log";

@Component({
    selector: 'app-log-detail',
    templateUrl: './log-detail.component.html',
    styleUrls: ['./log-detail.component.css']
})
export class LogDetailComponent implements OnInit {

    constructor() {
    }

    opened = false;
    item: ClusterLog = new ClusterLog();

    ngOnInit(): void {
    }

    open(item: ClusterLog) {
        item.message = item.message.replace(/[\\]/g, '');
        this.item = item;
        this.opened = true;
    }

}
