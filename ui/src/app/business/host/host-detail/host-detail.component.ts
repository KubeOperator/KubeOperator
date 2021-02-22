import {Component, OnInit} from '@angular/core';
import {Host} from '../host';


@Component({
    selector: 'app-host-detail',
    templateUrl: './host-detail.component.html',
    styleUrls: ['./host-detail.component.css']
})
export class HostDetailComponent implements OnInit {
    constructor() {
    }
    
    opened = false;
    item: Host = new Host();

    ngOnInit(): void {
    }

    onCancel() {
        this.item = new Host();
        this.opened = false;
    }

    open(item) {
        this.opened = true;
        this.item = item;
    }
}
