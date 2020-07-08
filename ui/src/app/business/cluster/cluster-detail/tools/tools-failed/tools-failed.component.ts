import {Component, OnInit} from '@angular/core';
import {ClusterTool} from "../tools";

@Component({
    selector: 'app-tools-failed',
    templateUrl: './tools-failed.component.html',
    styleUrls: ['./tools-failed.component.css']
})
export class ToolsFailedComponent implements OnInit {

    constructor() {
    }

    opened = false;
    item: ClusterTool = new ClusterTool();

    ngOnInit(): void {
    }

    open(item: ClusterTool) {
        this.opened = true;
        this.item = item;
    }

    onCancel() {
        this.opened = false;
    }
}
