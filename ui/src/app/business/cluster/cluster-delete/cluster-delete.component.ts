import {Component, OnInit} from '@angular/core';

@Component({
    selector: 'app-cluster-delete',
    templateUrl: './cluster-delete.component.html',
    styleUrls: ['./cluster-delete.component.css']
})
export class ClusterDeleteComponent implements OnInit {

    opened = false;

    constructor() {
    }


    ngOnInit(): void {
    }

    open() {
        this.opened = true;
    }

}
