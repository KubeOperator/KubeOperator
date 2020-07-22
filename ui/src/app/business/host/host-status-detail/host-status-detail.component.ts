import {Component, OnInit} from '@angular/core';
import {Host} from "../host";

@Component({
    selector: 'app-host-status-detail',
    templateUrl: './host-status-detail.component.html',
    styleUrls: ['./host-status-detail.component.css']
})
export class HostStatusDetailComponent implements OnInit {

    constructor() {
    }

    opened = false;
    item: Host = new Host();

    ngOnInit(): void {
    }

    open(item: Host) {
        this.item = item;
        this.opened = true;
    }

}
