import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ProjectResourceListComponent } from './project-resource-list.component';

describe('ProjectResourceListComponent', () => {
  let component: ProjectResourceListComponent;
  let fixture: ComponentFixture<ProjectResourceListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ProjectResourceListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProjectResourceListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
