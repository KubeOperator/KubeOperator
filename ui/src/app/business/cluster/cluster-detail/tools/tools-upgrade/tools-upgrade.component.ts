import {Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {ClusterTool} from '../tools';
import {Cluster} from '../../../cluster';
import {ToolsService} from '../tools.service';
import {TranslateService} from '@ngx-translate/core';

@Component({
    selector: 'app-tools-upgrade',
    templateUrl: './tools-upgrade.component.html',
    styleUrls: ['./tools-upgrade.component.css']
})
export class ToolsUpgradeComponent implements OnInit {

    constructor(private toolsService: ToolsService,
                private translateService: TranslateService,
    ) {}
    opened = false;
    isSubmitGoing = false;
    item: ClusterTool = new ClusterTool();
    @Input() currentCluster: Cluster;
    @Output() upgraded = new EventEmitter();


    ngOnInit(): void {
    }

    onSubmit() {
        this.isSubmitGoing = true;
        this.toolsService.upgrade(this.currentCluster.name, this.item).subscribe(data => {
            this.opened = false;
            this.upgraded.emit();
            this.isSubmitGoing = false;
        });
    }
    
    onCancel() {
        this.opened = false;
    }

    open(item: ClusterTool) {
        this.opened = true;
        this.item = item;
    }
}
