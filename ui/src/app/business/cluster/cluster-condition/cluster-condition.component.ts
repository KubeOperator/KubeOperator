import {Component, OnInit} from '@angular/core';

@Component({
    selector: 'app-cluster-condition',
    templateUrl: './cluster-condition.component.html',
    styleUrls: ['./cluster-condition.component.css']
})
export class ClusterConditionComponent implements OnInit {

    opened = false;
    item: string;

    constructor() {
    }

    ngOnInit(): void {
    }

    onCancel() {
        this.opened = false;
    }

    open(item: string) {
        this.item = item;
        this.opened = true;
    }

}
