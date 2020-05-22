import {Component, OnInit, ViewChild} from '@angular/core';
import {ClusterCreateRequest} from '../cluster';
import {NgForm} from '@angular/forms';

@Component({
    selector: 'app-cluster-create',
    templateUrl: './cluster-create.component.html',
    styleUrls: ['./cluster-create.component.css']
})
export class ClusterCreateComponent implements OnInit {

    opened = false;
    item: ClusterCreateRequest = new ClusterCreateRequest();

    @ViewChild('clusterForm') clusterForm: NgForm;

    constructor() {
    }

    ngOnInit(): void {
    }

    open() {
        this.item = new ClusterCreateRequest();
        this.clusterForm.resetForm();
        this.opened = true;
    }
}
