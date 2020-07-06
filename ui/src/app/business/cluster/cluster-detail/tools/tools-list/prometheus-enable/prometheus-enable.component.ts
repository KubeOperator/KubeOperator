import {Component, OnInit} from '@angular/core';
import {ClusterTool} from "../../tools";

@Component({
    selector: 'app-prometheus-enable',
    templateUrl: './prometheus-enable.component.html',
    styleUrls: ['./prometheus-enable.component.css']
})
export class PrometheusEnableComponent implements OnInit {

    constructor() {
    }

    opened = false;
    isSubmitGoing = false;
    item: ClusterTool;

    ngOnInit(): void {
    }

    open(item: ClusterTool) {
        this.opened = true;
        this.item = item;
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        this.opened = false;
    }

}
