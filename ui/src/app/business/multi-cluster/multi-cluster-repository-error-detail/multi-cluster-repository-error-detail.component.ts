import {Component, OnInit} from '@angular/core';

@Component({
    selector: 'app-multi-cluster-repository-error-detail',
    templateUrl: './multi-cluster-repository-error-detail.component.html',
    styleUrls: ['./multi-cluster-repository-error-detail.component.css']
})
export class MultiClusterRepositoryErrorDetailComponent implements OnInit {

    constructor() {
    }

    errorMessage: string;
    opened = false;

    ngOnInit(): void {
    }

    open(msg: string) {
        this.errorMessage = msg;
        this.opened = true;
    }


}
