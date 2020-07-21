import {Component, OnInit} from '@angular/core';
import {Project} from '../project';
import {ActivatedRoute, Router} from '@angular/router';

@Component({
    selector: 'app-project-detail',
    templateUrl: './project-detail.component.html',
    styleUrls: ['./project-detail.component.css']
})
export class ProjectDetailComponent implements OnInit {

    currentProject: Project;


    constructor(private router: Router, private route: ActivatedRoute) {
    }

    ngOnInit(): void {
        this.route.data.subscribe(data => {
            this.currentProject = data.project;
        });
    }

    backToProject() {
        this.router.navigate(['projects']);
    }
}
