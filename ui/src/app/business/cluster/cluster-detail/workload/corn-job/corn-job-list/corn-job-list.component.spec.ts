import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CornJobListComponent } from './corn-job-list.component';

describe('CornJobListComponent', () => {
  let component: CornJobListComponent;
  let fixture: ComponentFixture<CornJobListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CornJobListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CornJobListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
