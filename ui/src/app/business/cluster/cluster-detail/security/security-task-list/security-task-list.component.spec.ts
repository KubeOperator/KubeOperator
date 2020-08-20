import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SecurityTaskListComponent } from './security-task-list.component';

describe('SecurityTaskListComponent', () => {
  let component: SecurityTaskListComponent;
  let fixture: ComponentFixture<SecurityTaskListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SecurityTaskListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SecurityTaskListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
