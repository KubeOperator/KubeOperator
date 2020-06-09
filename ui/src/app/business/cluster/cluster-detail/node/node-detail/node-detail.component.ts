import {Component, OnInit} from '@angular/core';
import {V1Node} from '@kubernetes/client-node';

@Component({
    selector: 'app-node-detail',
    templateUrl: './node-detail.component.html',
    styleUrls: ['./node-detail.component.css']
})
export class NodeDetailComponent implements OnInit {

    opened = false;
    item: V1Node;

    constructor() {
    }

    ngOnInit(): void {
    }

    open(item: V1Node) {
        this.opened = true;
        this.item = item;
    }

}
