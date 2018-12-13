import { Component, OnInit } from '@angular/core';
import { OnDestroy } from "@angular/core";

import { NavigatorService } from "./navigator.service";
import { Nav } from "./navigator";

@Component({
  selector: 'app-navigator',
  templateUrl: './navigator.component.html',
  styleUrls: ['./navigator.component.css']
})
export class NavigatorComponent implements OnInit, OnDestroy {
  navList: Nav[] = [];
  showProjectsNav: boolean = true;
  showDetailNav: boolean = false;
  projectsUrl: string = '/ansible/projects';
  projectUrl: string;
  currentProjectName: '';
  // projectsNav: Nav[] = [
  //   {name: "Projects", shape: "organization", link: "/ansible/projects"},
  // ];
  // projectDetailNav: Nav[] = [
  //   {name: "Overview", shape: "dashboard", link: "overview"},
  //   {name: "Jobs", shape: "list", link: "projects"},
  //   {name: "Roles", shape: "grid-view", link: "roles"},
  //   {name: "AdHoc", shape: "code", link: "adhoc"},
  //   {name: "Inventory", shape: "vm", link: "get_inventory"},
  // ];
  subscription = null;

  constructor(private service: NavigatorService) {
    this.navList = [];
    this.subscription = this.service.messageObserver.subscribe(
      operation => {
        console.log(operation);
        if (operation.target === "detail") {
          this.showProjectsNav = false;
          this.showDetailNav = true;
          this.currentProjectName = operation.meta;
        } else {
          this.showProjectsNav = true;
          this.showDetailNav = false;
        }
      }
    )
  }

  ngOnInit() {
    console.log("Nav on init")
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
  }
}
