import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {CLusterImportRequest} from '../cluster';
import {NgForm} from '@angular/forms';
import {ClusterService} from '../cluster.service';
import {ActivatedRoute} from "@angular/router";
import {Project} from "../../project/project";

@Component({
    selector: 'app-cluster-import',
    templateUrl: './cluster-import.component.html',
    styleUrls: ['./cluster-import.component.css']
})
export class ClusterImportComponent implements OnInit {

    constructor(private clusterService: ClusterService, private route: ActivatedRoute) {
    }

    opened = false;
    item = new CLusterImportRequest();
    currentProject: Project;
    @Output() imported = new EventEmitter();
    @ViewChild('importForm') importForm: NgForm;

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentProject = data.project;
        });
    }

    open() {
        this.opened = true;
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        this.item.projectName = this.currentProject.name;
        this.clusterService.import(this.item).subscribe(() => {
            this.imported.emit();
            this.opened = false;
        });
    }

}
