import {Component, OnInit} from '@angular/core';
import {Cluster} from '../../cluster';

@Component({
    selector: 'app-backup',
    templateUrl: './backup.component.html',
    styleUrls: ['./backup.component.css']
})
export class BackupComponent implements OnInit {

    tab: string;
    currentCluster: Cluster;

    constructor() {
    }

    ngOnInit(): void {
        this.tab = 'strategy';
    }

    changeTab(tab: string) {
        this.tab = tab;
    }
}
