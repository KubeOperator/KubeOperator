import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { StorageClassListComponent } from './storage-class-list.component';

describe('StorageClassListComponent', () => {
  let component: StorageClassListComponent;
  let fixture: ComponentFixture<StorageClassListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ StorageClassListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(StorageClassListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
