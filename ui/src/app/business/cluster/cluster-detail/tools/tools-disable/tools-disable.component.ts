import {Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {NgForm} from '@angular/forms';
import {ClusterTool} from '../tools';
import {Cluster} from '../../../cluster';
import {ToolsService} from '../tools.service';

@Component({
    selector: 'app-tools-disable',
    templateUrl: './tools-disable.component.html',
    styleUrls: ['./tools-disable.component.css']
})
export class ToolsDisableComponent implements OnInit {

    constructor(private toolsService: ToolsService) {
    }

    opened = false;
    isSubmitGoing = false;
    @Output() disabled = new EventEmitter();
    tool: ClusterTool = new ClusterTool();
    @Input() currentCluster: Cluster;

    ngOnInit(): void {
    }

    open(item: ClusterTool) {
        this.tool = item;
        this.opened = true;
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        if (this.isSubmitGoing) {
            return;
        }
        this.isSubmitGoing = true;
        this.toolsService.disable(this.currentCluster.name, this.tool).subscribe(data => {
            this.isSubmitGoing = false;
            this.tool = data;
            this.opened = false;
            this.disabled.emit();
            location.reload();
        });
    }

}
