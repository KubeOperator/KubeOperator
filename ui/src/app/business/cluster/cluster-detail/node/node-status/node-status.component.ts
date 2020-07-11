import {Component, OnInit} from '@angular/core';
import {Node} from "../node";

@Component({
    selector: 'app-node-status',
    templateUrl: './node-status.component.html',
    styleUrls: ['./node-status.component.css']
})
export class NodeStatusComponent implements OnInit {

    constructor() {
    }

    opened = false;
    item: Node = new Node();

    ngOnInit(): void {
    }

    open(item: Node) {
        this.item = item;
        this.opened = true;
    }

    onCancel() {
        this.opened = false;
    }

}
