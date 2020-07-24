import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ProjectMemberCreateComponent } from './project-member-create.component';

describe('ProjectMemberCreateComponent', () => {
  let component: ProjectMemberCreateComponent;
  let fixture: ComponentFixture<ProjectMemberCreateComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ProjectMemberCreateComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProjectMemberCreateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
